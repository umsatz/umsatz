require.config({

  /* starting point for application */
  deps: [
    'backbone',
    'backbone.marionette',

    'ractive',

    'components/currency',
    'components/account-picker',

    'i18n',
    'main'
  ],
  waitSeconds: 60,

  shim: {
    underscore: {
      exports: '_'
    },
    backbone: {
      deps: ['underscore', 'jquery'],
      exports: 'Backbone'
    },
    'backbone.marionette': {
      deps: ['backbone'],
      exports: 'Backbone'
    },
    'backbone.relational': {
      deps: ['backbone'],
      exports: 'Backbone'
    },
    accounting: {
      exports: 'accounting'
    },
    foundation: {
      deps: ['jquery'],
      exports: 'Foundation'
    },
    'foundation.topbar': {
      deps: ['foundation'],
      exports: 'Foundation'
    },
    'foundation.dropdown': {
      deps: ['foundation'],
      exports: 'Foundation'
    },
    typeahead: {
      deps: ['jquery', 'bloodhound']
    },
    bloodhound: {
      deps: ['jquery'],
      exports: 'Bloodhound'
    },
    pdfviewer: {
      deps: ['pdfjs']
    }
  },

  paths: {
    jquery: '../bower_components/jquery/dist/jquery',
    backbone: '../bower_components/backbone/backbone',
    underscore: '../bower_components/underscore-amd/underscore',

    /* alias all marionette libs */
    'backbone.marionette': '../bower_components/backbone.marionette/lib/backbone.marionette',
    'backbone.radio': '../bower_components/backbone.radio/build/backbone.radio',
    'backbone.cacheit': 'vendor/backbone.cacheit',

    /* ractive */
    ractive: '../bower_components/ractive/ractive',
    'ractive-backbone': '../bower_components/ractive-adaptors-backbone/ractive-adaptors-backbone',
    'amd-loader': '../bower_components/requirejs-ractive/amd-loader',
    'pdfjs': '../bower_components/pdfjs-dist/build/pdf',
    'pdfviewer': '../bower_components/pdfjs-dist/web/pdf_viewer',
    rv: '../bower_components/requirejs-ractive/rv',
    rvc: '../bower_components/requirejs-ractive/rvc',

    /* i18n */
    i18next: '../bower_components/i18next/i18next.amd',

    /* foundation */
    foundation: '../bower_components/foundation/js/foundation/foundation',
    'foundation.topbar': '../bower_components/foundation/js/foundation/foundation.topbar',
    'foundation.dropdown': '../bower_components/foundation/js/foundation/foundation.dropdown',

    /* auto completion */
    typeahead: '../bower_components/typeahead.js/dist/typeahead.bundle',
    bloodhound: '../bower_components/typeahead.js/dist/bloodhound',

    /* money */
    accounting: '../bower_components/accounting/accounting',

    /* Alias text.js for template loading and shortcut the templates dir to tmpl */
    text: '../bower_components/requirejs-text/text',
    tmpl: '../templates',
  }
});
