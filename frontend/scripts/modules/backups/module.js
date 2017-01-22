'use strict';
define([
  'application',
  'communicator',
  './collections/backups',
  './routers/backups',

], function(App, Communicator, BackupCollection, Router) {
  var Backups = App.module('Backups');

  var backups = new BackupCollection();

  Communicator.reqres.reply('backups', function() {
    return backups;
  });

  App.addInitializer(function() {
    backups.fetch();

    new Router();
  });

  return Backups;
});
