define([
  'application',
  'rv!../templates/signup',
  'communicator'
], function(App, signupTemplate, comm) {
  'use strict';

  return function() {
    var ractive = new App.Ractive({
      template: signupTemplate,
      adapt: ['Backbone'],

      el: '#content'
    });

    ractive.on('signup', function(context) {
      context.original.preventDefault();

      comm.reqres.request('session:validate', ractive.data.signup.password, ractive.data.signup.totp)
        .then(function() {
          App.Backbone.history.navigate('/', true);
        }).fail(function() {
          // TODO display errors
        });
    });

    ractive.on('validate', function(context) {
      context.original.preventDefault();

      jQuery.post('/api/auth/signup/', JSON.stringify({
        firstname: context.context.signup.firstname,
        lastname:  context.context.signup.lastname,
        company:   context.context.signup.company,
        email:     context.context.signup.email,
        password:  context.context.signup.password,
        tax_id:    context.context.signup.tax_id
      })).then(function(rawResp) {
        var resp = JSON.parse(rawResp);
        var qr = {
          href: window.location.protocol + '//' + window.location.hostname + '/api/auth/' + resp.qrcode_url
        };
        ractive.set('qr', qr);
      }).fail(function() {
        console.log('error', arguments);
        ractive.set('qr', null);
        // TODO display error
      });
    });

    return ractive;
  };

});
