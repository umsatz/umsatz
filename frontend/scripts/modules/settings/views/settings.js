define([
  'application',
  'communicator',
  'rv!../templates/settings'
], function(App, comm, template) {
  'use strict';

  return function(backups) {
    var ractive = new App.Ractive({
      template: template,
      adapt: ['Backbone'],

      el: 'content',

      data: {
        backups: backups
      }
    });

    comm.reqres.request('registration').then(function(req) {
      ractive.set('settings', req);
    });

    ractive.on('save', function(ctx) {
      ctx.original.preventDefault();

      jQuery.ajax({
        url: '/api/auth/registration/',
        data: JSON.stringify({
          firstname: ctx.context.settings.firstname,
          lastname:  ctx.context.settings.lastname,
          company:   ctx.context.settings.company,
          email:     ctx.context.settings.email,
          tax_id:    ctx.context.settings.tax_id
        }),
        type: 'put'
      }).then(function(rawResp) {
        console.log('ok', rawResp);
      });
    });

    return ractive;
  };
});
