define([
  'backbone.marionette',
  '../views/index'
], function(Marionette, backupsOverview) {
  'use strict';

  return Marionette.AppRouter.extend({

    controller: {
      backupsOverview: function() {
        var comm = require('communicator');
        backupsOverview(comm.reqres.request('backups'));
      }
    },

    appRoutes: {
      'backups': 'backupsOverview',
    },

  });
});
