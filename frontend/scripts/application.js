define([
  'backbone',
  'backbone.marionette',
  'backbone.cacheit',
  'i18n',
  'ractive',
  'ractive-backbone',
  'components/currency',
], function(Backbone, Marionette, _, I18n, Ractive) {
  'use strict';
  var app = new Backbone.Marionette.Application();

  app.I18n = I18n;
  app.Backbone = Backbone;
  app.Ractive = Ractive.extend({
    data: {
      t: I18n.t
    }
  });

  return app;
});
