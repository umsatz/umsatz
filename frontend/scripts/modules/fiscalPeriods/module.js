'use strict';
define([
  'application',
  'communicator',
  'collections/fiscalPeriods',
  './routers/fiscalPeriods',
  'components/account-picker'
], function(App, Communicator, PositionCollection, Router) {
  var FiscalPeriod = App.module('FiscalPeriods');

  var positions = new PositionCollection();

  Communicator.reqres.reply('fiscalPeriods', function() {
    return positions;
  });

  Communicator.reqres.reply('fiscalPeriods:get', function(id) {
    var fiscalPeriodId = parseInt(id, 10),
        promise = new $.Deferred();

    positions.fetch()
      .then(function() {
        promise.resolve(positions.get(fiscalPeriodId));
      })
      .fail(promise.reject);

    return promise;
  });

  App.addInitializer(function() {
    positions.fetch();

    new Router();
  });

  return FiscalPeriod;
});
