define([
  'application',
  'jquery',
  'backbone',
  'communicator',
  '../collections/positions',
  '../models/position',
  '../views/index',
  '../views/edit'
], function(App, $, Backbone, Communicator, PositionCollection, Position, positionsOverview, positionForm) {
  'use strict';

  var activeContent = null;
  var comm = require('communicator');
  var positions = new PositionCollection();
  var fiscalYear = null;

  var redirectToIndex = function redirectToIndex() {
    var url = Communicator.reqres.request('positions:lastUrl');
    if (url === '') {
      url = 'fiscalPeriods/' + fiscalYear.get('id');
    }
    Backbone.history.navigate(url, true);
  };

  var updateFiscalPeriod = function updateFiscalPeriod(period, positions) {
    var totalIncomeCents = 0,
      totalExpenseCents = 0;

    positions.forEach(function(position) {
      if (position.isIncome()) {
        totalIncomeCents += position.get('totalAmountCentsEur');
      } else {
        totalExpenseCents += position.get('totalAmountCentsEur');
      }
    });

    period.set({
      totalIncomeCents: totalIncomeCents,
      totalExpenseCents: totalExpenseCents,
      positionsCount: positions.length
    });
  };

  return App.Backbone.Marionette.AppRouter.extend({

    controller: {
      loadFiscalYear: function(fiscalPeriodId) {
        var promise = new $.Deferred();
        fiscalPeriodId = parseInt(fiscalPeriodId, 10);

        comm.reqres.request('fiscalPeriods:get', fiscalPeriodId).done(function(fiscalPeriod) {
          if (fiscalYear && fiscalYear.get('id') !== fiscalPeriodId) {
            positions = new PositionCollection();
          }

          fiscalYear = fiscalPeriod;

          positions.url = fiscalYear.get('positionsUrl');
          positions.fetch({ reload: true }).done(function() {
            positions.sort();
            promise.resolve();
          });

          if (activeContent !== null) {
            activeContent.teardown();
          }
        })
        .fail(function() {
          promise.reject();
        });

        return promise;
      }.bind(this),

      index: function(fiscalPeriodId, params) {
        this.loadFiscalYear(fiscalPeriodId).done(function() {
          activeContent = positionsOverview(fiscalYear, positions, params || {});

          activeContent.on('delete', function(event) {
            event.context.destroy();
            updateFiscalPeriod(fiscalYear, positions);
          });
        });
      },

      indexWithSearch: function(fiscalPeriodId, search) {
        this.index(fiscalPeriodId, {
          search: search
        });
      },

      edit: function(fiscalPeriodId, id) {
        this.loadFiscalYear(fiscalPeriodId).done(function() {
          var position = positions.get(parseInt(id, 10));
          position.url = '/api/positions/' + id;
          activeContent = positionForm(fiscalYear, position);

          activeContent.on('fiscalItem:put', function() {
            updateFiscalPeriod(fiscalYear, positions);
            redirectToIndex();
          });
          activeContent.on('fiscalItem:cancel', redirectToIndex);
        });
      },

      new: function(fiscalPeriodId, positionTemplate) {
        this.loadFiscalYear(fiscalPeriodId).done(function() {
          var position = null;
          if (positionTemplate === undefined || positionTemplate === null) {
            position = new Position({
              fiscalPeriodId: fiscalYear.get('id')
            });
          } else {
            position = positionTemplate.clone();
            position.unset('id');
          }
          activeContent = positionForm(fiscalYear, position);

          activeContent.on('fiscalItem:put', function() {
            positions.add(position);
            updateFiscalPeriod(fiscalYear, positions);
            redirectToIndex();
          }.bind(this));

          activeContent.on('fiscalItem:cancel', redirectToIndex);
        });
      },

      clone: function(fiscalPeriodId, positionId) {
        this.loadFiscalYear(fiscalPeriodId).done(function() {
          this.new(fiscalPeriodId, positions.get(parseInt(positionId, 10)));
        }.bind(this));
      }
    },

    appRoutes: {
      'fiscalPeriods/:fiscal_period_id/s/:search': 'indexWithSearch',
      'fiscalPeriods/:fiscal_period_id/positions/new': 'new',
      'fiscalPeriods/:fiscal_period_id/positions/:id/edit': 'edit',
      'fiscalPeriods/:fiscal_period_id/positions/:id/clone': 'clone',
    },

    initialize: function() {
      this.appRoute(/fiscalPeriods\/(\d+)$/, 'index');
    }
  });
});
