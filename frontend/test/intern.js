// Learn more about configuring this file at <https://github.com/theintern/intern/wiki/Configuring-Intern>.
// These default settings work OK for most people. The options that *must* be changed below are the
// packages, suites, excludeInstrumentation, and (if you want functional tests) functionalSuites.
/* global define */
define({
  // The port on which the instrumenting proxy will listen
  proxyPort: 9090,

  // A fully qualified URL to the Intern proxy
  proxyUrl: 'http://localhost:9090/',

  reporters: ['console', 'lcovhtml'],

  webdriver: {
    host: 'localhost',
    port: 4444
  },

  // Default desired capabilities for all environments. Individual capabilities can be overridden by any of the
  // specified browser environments in the `environments` array below as well. See
  // https://code.google.com/p/selenium/wiki/DesiredCapabilities for standard Selenium capabilities and
  // https://saucelabs.com/docs/additional-config#desired-capabilities for Sauce Labs capabilities.
  // Note that the `build` capability will be filled in with the current commit ID from the Travis CI environment
  // automatically
  capabilities: {
    'selenium-version': '2.42.0'
  },

  // Browsers to run integration testing against. Note that version numbers must be strings if used with Sauce
  // OnDemand. Options that will be permutated are browserName, version, platform, and platformVersion; any other
  // capabilities options specified for an environment will be copied as-is
  environments: [
    // {
    //   browserName: 'firefox'
    // },
    {
      browserName: 'chrome'
    },
    // { browserName: 'safari' },
  ],

  // Maximum number of simultaneous integration tests that should be executed on the remote WebDriver service
  maxConcurrency: 3,

  // Name of the tunnel class to use for WebDriver tests
  // tunnel: 'SauceLabsTunnel',

  // The desired AMD loader to use when running unit tests (client.html/client.js). Omit to use the default Dojo
  // loader
  // useLoader: {
  // 'host-node': 'dojo/dojo',
  // 'host-browser': '../bower_components/requirejs/require.js'
  // },

  // Configuration options for the module loader; any AMD configuration options supported by the specified AMD loader
  // can be used here
  loader: {
    // Packages that should be registered with the loader in each testing environment
    packages: [
      {
        name: 'test',
        location: 'test'
      },
      {
        name: 'scripts',
        location: 'scripts'
      },
      {
        name: 'accounting',
        main: 'accounting',
        location: 'bower_components/accounting'
      },
      {
        name: 'sinon',
        location: 'node_modules/sinon/pkg'
      },
      {
        name: 'backbone',
        location: 'bower_components/backbone',
        main: 'backbone'
      },
      {
        name: 'underscore',
        location: 'bower_components/underscore',
        main: 'underscore'
      },
      {
        name: 'jquery',
        location: 'bower_components/jquery/dist',
        main: 'jquery'
      }
    ]
  },

  // Non-functional test suite(s) to run in each browser
  suites: [
    'test/template-helpers/currencySpec',
    'test/models/accountSpec',
    'test/models/fiscalPeriodSpec',
  ],

  // Functional test suite(s) to run in each browser once non-functional tests are completed
  functionalSuites: [],

  // A regular expression matching URLs to files that should not be included in code coverage analysis
  excludeInstrumentation: /^(?:test|bower_components|node_modules)\//
});