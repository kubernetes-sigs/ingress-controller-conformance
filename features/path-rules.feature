@extensions/v1beta1
Feature: Ingress Path Rules
  I need to be able to do HTTP calls to an Ingress with path rules.

  Scenario: Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service
  (/foo matches /foo)
    Given I have an Ingress named "path-rules" in the "default" namespace
    When I send a "GET" "http://path-rules/foo" request
    Then the response status-code must be 200
    And the response must be served by the "path-rules-foo" service
    #And the request path must be "/foo"

  Scenario: Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service
  (/foo matches /foo/)
    Given I have an Ingress named "path-rules" in the "default" namespace
    When I send a "GET" "http://path-rules/foo/" request
    Then the response status-code must be 200
    And the response must be served by the "path-rules-foo" service
    #And the request path must be "/foo/"