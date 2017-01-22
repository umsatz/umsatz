define([
  'backbone.marionette',
  'communicator',
  '../views/preview',
  '../models/template'
], function(Marionette, comm, previewTemplate, Template) {
  'use strict';

  return Marionette.AppRouter.extend({

    controller: {
      preview: function(fiscalPeriodId) {
        $.when(
          comm.reqres.request('fiscalPeriods:get', fiscalPeriodId),
          new Template({ id: fiscalPeriodId }).fetch()
        ).then(function(fiscalPeriod, template) {
          previewTemplate(fiscalPeriod, template);
        });
      }
    },

    appRoutes: {
      'fiscalPeriods/:id/export': 'preview'
    },

  });
});
