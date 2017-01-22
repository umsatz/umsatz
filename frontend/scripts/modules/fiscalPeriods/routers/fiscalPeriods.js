define([
  'backbone',
  'backbone.marionette',
  'communicator',
  'models/fiscalPeriod',
  '../views/index',
  '../views/edit'
], function(Backbone, Marionette, Communicator, FiscalPeriod, fiscalOverviewView, periodForm) {
  'use strict';

  var activeContent = null;

  return Marionette.AppRouter.extend({
    controller: {
      overview: function() {
        fiscalOverviewView(Communicator.reqres.request('fiscalPeriods'));
      },

      new: function() {
        var period = new FiscalPeriod();
        activeContent = periodForm(period);

        activeContent.on('fiscalPeriod:put', function(period) {
          Communicator.reqres.request('fiscalPeriods').add(period);
          Backbone.history.navigate('/', true);
        }.bind(this));

        activeContent.on('fiscalPeriod:cancel', function() {
          Backbone.history.navigate('/', true);
        }.bind(this));
      },

      edit: function(id) {
        var period = Communicator.reqres.request('fiscalPeriods').get(id);

        activeContent = periodForm(period);
        activeContent.on('fiscalPeriod:put', function() {
          Backbone.history.navigate('/', true);
        }.bind(this));

        activeContent.on('fiscalPeriod:cancel', function() {
          Backbone.history.navigate('/', true);
        }.bind(this));
      }
    },

    appRoutes: {
      '': 'overview',
      'fiscalPeriods/new': 'new',
      'fiscalPeriods/:id/edit': 'edit'
    },

  });
});
