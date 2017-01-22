define([
  'backbone',
  '../models/backup'
], function(Backbone, Backup) {
  'use strict';

  return Backbone.Collection.extend({
    model: Backup,
    url: '/api/backups'
  });
});
