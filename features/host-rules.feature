@extensions/v1beta1
Feature: Ingress Host Rules
  I need to be able to do HTTP calls to an Ingress with host rules.

  Scenario: Ingress with host rule should send traffic to the correct backend service
    Given I have an Ingress named "host-rules" in the "default" namespace
    When I send a "GET" "http://foo.bar.com" request
    Then the response status-code must be 200
    And the response must be served by the "host-rules" service
    And the request host must be "foo.bar.com"