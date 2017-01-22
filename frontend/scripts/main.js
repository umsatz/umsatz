/* global Backbone */
require([
  'jquery',
  'application',
  'communicator',
  'collections/accounts',
], function($, App, Communicator, Accounts) {
  'use strict';

  $.ajaxPrefilter('*', function(options, originalOptions, jqXHR) {
    if (Communicator.reqres.request('session:current') != null) {
      jqXHR.setRequestHeader('X-UMSATZ-SESSION', Communicator.reqres.request('session:current').key);
    }
  });

  var accounts = new Accounts();
  Communicator.reqres.reply('accounts', function() {
    return accounts;
  });

  App.getCurrentRoute = function(){
    return Backbone.history.fragment;
  };

  App.loadDashboard = function() {
    require(['views/navigation', 'views/footer']);
    accounts.fetch();

    Communicator.reqres.request('fiscalPeriods').fetch({ reload: true });
  };

  App.listenTo(Communicator, 'user:signin', App.loadDashboard.bind(this));

  App.on('start', function() {
    if (Backbone.history) {
      require([
        'modules/backups/module',
        'modules/exports/module',
        'modules/fiscalPeriods/module',
        'modules/positions/module',
        'modules/login/module',
        'modules/settings/module',
      ], function () {
        Backbone.history.start();

        $.get('/api/').fail(function(resp) {
          if (resp.status === 401) {
            Backbone.history.navigate('/login', true);
          }
        }).then(App.loadDashboard.bind(this));

      });
    }
  });

  App.start();
});
