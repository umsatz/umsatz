'use strict';
define([
  'application',
  './routers/exports',
], function(App, Router) {
  var Exports = App.module('Exports');

  App.addInitializer(function() {
    new Router();
  });

  return Exports;
});
