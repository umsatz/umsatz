define([
  'application',
  'communicator',
  'helpers/extend',
  '../helpers/summation',
  'rv!../templates/index'
], function(App, Communicator, extend, sumMixin, PositionsTemplate) {
  'use strict';

  return function(fiscalYear, positions, params) {
    var MonthComponent = App.Ractive.extend({
      template: '{{ month }}',
      init: function() {
        this.observe('month', function(value) {
          if (value !== '' && value != null) {
            this.set({
              month: App.I18n.t('date.month_abbrs.' + parseInt(value, 10))
            });
          }
        });
      },
      lazy: true,
      twoway: false
    });

    var filtered = positions.select(function() {
      return true;
    });
    extend(filtered, sumMixin);
    var months = [];

    function updateMonths() {
      months = [];
      for (var i = 0; i < filtered.length; i++) {
        var pos = filtered[i];
        if (i === 0) {
          months[i] = pos.invoiceMonth();
          continue;
        }

        var prev = filtered[i - 1];
        if (prev.invoiceMonth() !== pos.invoiceMonth()) {
          months[i] = pos.invoiceMonth();
          continue;
        }

        months[i] = null;
      }

      ractive.data.months = months;
      ractive.update('months');
    }

    function updateFiltered(query) {
      if (query === undefined) {
        query = '';
      }

      if (query === '') {
        ractive.data.filtered = filtered = positions.select(function() {
          return true;
        });
      } else {
        var regexp = new RegExp(query.toString().toLowerCase());
        ractive.data.filtered = filtered = positions.select(function(position) {
          return regexp.test(position.get('description').toLowerCase()) || regexp.test(position.get('invoiceNumber').toLowerCase());
        });
      }

      extend(filtered, sumMixin);
      var url = '/fiscalPeriods/' + fiscalYear.get('id');
      if (query !== '') {
        url += '/s/' + query;
      }
      Communicator.reqres.request('positions:lastUrl', url);
      App.Backbone.history.navigate(url, {
        replace: true,
        trigger: false
      });

      ractive.update('filtered');
    }

    var ractive = new App.Ractive({
      template: PositionsTemplate,
      adapt: ['Backbone'],

      el: 'content',

      data: {
        months: months,
        fiscalYear: fiscalYear,
        positions: positions,
        filtered: filtered,
        searchQuery: params.search || '',
        shortDate: function(content) {
          var parts = content.split('-');

          if (parts.length === 3) {
            return [parts[2], parts[1], ''].join('.');
          }
          return '';
        },
        vat: function(amountCents, tax) {
          return amountCents - (amountCents / (tax / 100 + 1));
        }
      },

      components: {
        month: MonthComponent
      }
    });
    updateMonths();

    ractive.observe('searchQuery', function(query) {
      updateFiltered(query);
      updateMonths();
    });

    positions.on('remove', function() {
      updateFiltered(ractive.data.searchQuery);
      updateMonths();
    });

    ractive.on('clearSearch', function(event) {
      event.original.preventDefault();
      event.original.stopPropagation();

      this.set({
        searchQuery: ''
      });
    });

    return ractive;
  };
});
