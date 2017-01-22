define([
  'backbone',
  'backbone.marionette',
  'backbone.radio'
], function(Backbone) {
  'use strict';

  return _.extend({
    reqres: Backbone.Radio.Requests,
    commands: Backbone.Radio.Commands
  }, Backbone.Events);
});
