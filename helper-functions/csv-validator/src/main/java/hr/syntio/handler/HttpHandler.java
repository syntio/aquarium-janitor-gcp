
// Copyright 2020 Syntio Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package serves as a medium for invoking Cloud Functions for csv and xml validation
package hr.syntio.handler;

import com.google.cloud.functions.HttpFunction;
import com.google.cloud.functions.HttpRequest;
import com.google.cloud.functions.HttpResponse;
import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;
import hr.syntio.csv.Validator;

import java.io.IOException;
import java.net.HttpURLConnection;

public class HttpHandler implements HttpFunction {

    private static final Gson gson = new Gson();

    public static class RequestStructure {
        RequestStructure(String data, String schema) {
            this.data = data;
            this.schema = schema;
        }

        String data;
        String schema;

        boolean isValidStructure() {
            return this.data != null && this.schema != null;
        }
    }

    public static class ResponseStructure {
        ResponseStructure(boolean validation, String info) {
            this.validation = validation;
            this.info = info;
        }

        boolean validation;
        String info;
    }

    @Override
    public void service(HttpRequest request, HttpResponse response) throws IOException {
        String data, schema;
        String contentType = request.getContentType().orElse("");

        switch (contentType) {
            case "application/json":
                try {
                    RequestStructure requestStructure = gson.fromJson(request.getReader(), RequestStructure.class);
                    if (requestStructure.isValidStructure()) {
                        data = requestStructure.data;
                        schema = requestStructure.schema;
                    } else {
                        handleResponse(response, new ResponseStructure(false, "Request json body doesn't have 'data' and 'schema' fields."), HttpURLConnection.HTTP_BAD_REQUEST);
                        return;
                    }
                } catch (JsonSyntaxException exception) {
                    handleResponse(response, new ResponseStructure(false, "Request body isn't a valid json."), HttpURLConnection.HTTP_BAD_REQUEST);
                    return;
                } catch (Exception exception) {
                    handleResponse(response, new ResponseStructure(false, "Request deserialization problem."), HttpURLConnection.HTTP_BAD_REQUEST);
                    return;
                }
                break;
            default:
                handleResponse(response, new ResponseStructure(false, "Request content type isn't 'application/json'."), HttpURLConnection.HTTP_BAD_REQUEST);
                return;
        }

        boolean isValid = Validator.validate(data, schema);
        handleResponse(response, new ResponseStructure(isValid, "successful validation"), HttpURLConnection.HTTP_OK);
    }

    private void handleResponse(HttpResponse response, ResponseStructure responseStructure, int statusCode) throws IOException {
        String responseString = gson.toJson(responseStructure);
        response.setContentType("application/json");
        response.setStatusCode(statusCode);
        response.getWriter().write(responseString);
        response.getWriter().close();
    }
}

