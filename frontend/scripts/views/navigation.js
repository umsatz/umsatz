define([
  'application',
  'rv!templates/navigation',
  'jquery',
  'communicator',
  'foundation.topbar'
], function(App, NavigationTemplate, jQuery, comm) {
  'use strict';

  var fiscalPeriods = comm.reqres.request('fiscalPeriods');

  var navigation = new App.Ractive({
    template: NavigationTemplate,
    adapt: ['Backbone'],

    el: '[role=navigation]',

    data: {
      fiscalPeriods: fiscalPeriods.select(function(fiscalPeriod) {
        return !fiscalPeriod.get('archived');
      })
    },

    complete: function() {
      jQuery(document).foundation();
    }
  });

  fiscalPeriods.on('sync', function() {
    navigation.set('fiscalPeriods', fiscalPeriods.select(function(fiscalPeriod) {
      return !fiscalPeriod.get('archived');
    }));
  });

  return navigation;
});
