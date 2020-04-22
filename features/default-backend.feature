@networking.k8s.io/v1beta1 @extensions/v1beta1
Feature: Ingress Default Backend
  I need to be able to do HTTP calls to an Ingress with a single default backend and no rules.

  Scenario: Ingress with no rules should send traffic to the correct backend service
    Given I have an Ingress named "single-service" in the "default" namespace
    When I send a "GET" "http://any-host" request
    Then the response status-code must be 200
    And the response must be served by the "single-service" service
    And the response proto must be "HTTP/1.1"
    And the response headers must contain <key> with matching <value>
        | key            | value |
        | Content-Length | *     |
        | Content-Type   | *     |
        | Date           | *     |
        | Server         | *     |
    And the request method must be "GET"
    And the request proto must be "HTTP/1.1"
    And the request headers must contain <key> with matching <value>
        | key        | value              |
        | User-Agent | Go-http-client/1.1 |