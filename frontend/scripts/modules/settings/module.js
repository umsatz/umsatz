'use strict';
define([
  'application',
  'communicator',
  './views/settings'
], function(App, Communicator, Settings) {
  var Module = App.module('Settings');

  var Router = App.Backbone.Marionette.AppRouter.extend({
    controller: {
      settings: function() {
        Settings();
      }
    },

    appRoutes: {
      'settings': 'settings'
    }
  });

  App.addInitializer(function() {
    new Router();
  });

  return Module;
});
