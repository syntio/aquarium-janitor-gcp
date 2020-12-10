# Copyright 2020 Syntio Inc.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Package handles xml schema validation
import json
import xmlschema
from flask import Response


def http_validation_handler(request):
    request_json = request.get_json(silent=True)
    is_valid = False

    if request_json and 'data' in request_json and 'schema' in request_json:
        data = request_json['data']
        schema = request_json['schema']
        try:
            is_valid = validate(data, schema)
            response = make_response(is_valid, "successful validation", 200)
        except:
            response = make_response(is_valid, "invalid json body content, can't resolve 'data' and 'schema' fields.", 400)
    else:
        response = make_response(False, "invalid request, needs 'data' and 'schema' fields.", 400)

    return response


def validate(data, schema):
    schema = xmlschema.XMLSchema(schema)
    return schema.is_valid(data)


def make_response(validation, info, status):
    response_data = {
        "validation": validation,
        "info": info
    }

    response = Response()
    response.data = json.dumps(response_data)
    response.status_code = status
    return response
