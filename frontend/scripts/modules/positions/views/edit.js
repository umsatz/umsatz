define([
  'application',
  'rv!../templates/edit',
  'helpers/computedPropertyWrapper',
  'modules/preview/module'
], function(App, PositionsTemplate, propertyWrapper, Preview) {
  'use strict';

  var accountLabel = function(obj) {
    return propertyWrapper('label', obj);
  };
  var uploadFile = function(fileList, prefix) {
    // clear the filename of unwanted characters
    var clean = function(string) {
      return string.replace(/\s/, '_');
    };

    var promise = new $.Deferred();
    if (fileList !== null) {
      var file = fileList[0];
      $.ajax({
        url: '//' + window.location.host + '/upload/' + prefix + '/' + clean(file.name),
        method: 'POST',
        data: file,
        processData: false
      }).done(function(data) {
        var downloadPath = '/download/' + data.file.key;
        promise.resolve(downloadPath);
      }).fail(function() {
        promise.reject();
      });
    } else {
      promise.resolve();
    }
    return promise;
  };

  return function(fiscalPeriod, position) {
    var comm = require('communicator');
    var accounts = comm.reqres.request('accounts');

    var originalAttributes = _.clone(position.attributes);

    var ractive = new App.Ractive({
      template: PositionsTemplate,
      adapt: ['Backbone'],

      el: 'content',

      data: {
        position: position,
        year: fiscalPeriod.get('year'),
        convertedTotalAmountCents: 0,
        accountLabelFrom: '',
        accountLabelTo: ''
      },

      computed: {
        previewable: {
          get: function() {
            return Preview.canPreview(position.get('attachmentPath'));
          }
        },
        totalAmount: {
          get: function() {
            return (this.get('position.totalAmountCents') / 100.0).toFixed(2);
          },
          set: function(value) {
            this.set({
              'position.totalAmountCents': parseInt(value * 100)
            });
          }
        }
      }
    });

    accounts.fetch().then(function() {
      ractive.set('accountLabelFrom', ( position.get('accountCodeFrom') !== '' ? accounts.findWhere({
        code: position.get('accountCodeFrom')
      }).get('label') : ''));
      ractive.set('accountLabelTo', ( position.get('accountCodeTo') !== '' ? accounts.findWhere({
        code: position.get('accountCodeTo')
      }).get('label') : ''));
    });

    var refreshExchangeInfos = function() {
      if (position.get('currency') === 'EUR') {
        return;
      }

      var promise = new $.Deferred();
      $.get('/api/rates/' + position.get('invoiceDate'))
      .done(function(exchangeInfos) {
        var rate = exchangeInfos.rates[position.get('currency')];
        ractive.set('exchangeRate', (1 / rate).toFixed(4));
        ractive.set('convertedTotalAmountCents', position.get('totalAmountCents') / rate);
      })
      .fail(function() {
        promise.fail();
      });
      return promise;
    };
    ractive.observe('position.invoiceDate', refreshExchangeInfos.bind(this));
    ractive.observe('position.currency', refreshExchangeInfos.bind(this));
    ractive.observe('position.totalAmountCents', refreshExchangeInfos.bind(this));

    if (Preview.canPreview(position.get('attachmentPath'))) {
      Preview.display($('#viewerContainer'), position.get('attachmentPath'));
    }

    ractive.on({
      save: function(event) {
        event.original.preventDefault();
        var model = event.context.position;

        // TODO maybe handle account update errors?
        accounts.upsert({
          label: ractive.data.accountLabelFrom,
          code: model.get('accountCodeFrom').toString()
        });
        accounts.upsert({
          code: model.get('accountCodeTo').toString(),
          label: ractive.data.accountLabelTo
        });

        uploadFile(position.get('attachment'), fiscalPeriod.get('id'))
        .done(function(attachmentPath) {
          if (attachmentPath !== null) {
            model.set('attachmentPath', attachmentPath);
          }
          model.save()
          .then(function() {
            model.set({
              errors: {}
            });
            this.fire('fiscalItem:put');
          }.bind(this))
          .fail(function(response) {
            position.set({
              errors: response.responseJSON.errors
            });
          }.bind(this));
        }.bind(this)
        );
      },
      teardown: function() {
        if (position.isNew()) { // we probably used the navigation to close this view
          position.destroy();
        }
      },
      cancel: function() {
        position.set(originalAttributes, {
          silence: true
        });
        this.fire('fiscalItem:cancel');
      }
    });

    return ractive;
  };
});
