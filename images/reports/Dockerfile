# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# build stage
FROM maven:3.6.3-openjdk-14 as builder

WORKDIR /build

COPY pom.xml /build/
RUN mvn dependency:go-offline

COPY src /build/src

RUN mvn clean install

FROM openjdk:14.0.2-slim-buster

WORKDIR /app

COPY --from=builder /build/target/*jar-with-dependencies.jar report-builder.jar

# https://stackoverflow.com/questions/58991966/what-java-security-egd-option-is-for
ENV JAVA_OPTS "-Djava.security.egd=file:/dev/./urandom"

CMD ["/bin/bash", "-c", "java $JAVA_OPTS -jar report-builder.jar"]
