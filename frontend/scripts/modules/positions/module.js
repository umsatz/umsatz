'use strict';
define([
  'application',
  'communicator',
  './collections/positions',
  './routers/positions',
], function(App, Communicator, PositionCollection, Router) {
  var FiscalPeriod = App.module('FiscalPeriodPositions');

  FiscalPeriod.positions = new PositionCollection();

  var positionUrl = '';

  Communicator.reqres.reply('positions:lastUrl', function(url) {
    if (url !== undefined) {
      positionUrl = url;
    }
    return positionUrl;
  });

  App.addInitializer(function() {
    new Router();
  });

  return FiscalPeriod;
});
