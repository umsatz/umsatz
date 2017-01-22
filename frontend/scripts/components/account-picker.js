define(
[
  'ractive',
  'bloodhound',
  'typeahead'
], function(Ractive, Bloodhound) {
  'use strict';

  var setupAutocompletion = function setupAutocompletion() {
    var comm = require('communicator');
    var accounts = comm.reqres.request('accounts');

    this.datasource = new Bloodhound({
      datumTokenizer: function(d) {
        return _.map([d.code, d.label].concat(d.label.split(',')), $.trim);
      },
      queryTokenizer: Bloodhound.tokenizers.whitespace,
      local: function() {
        return accounts.toJSON();
      }
    });
    this.datasource.initialize();

    var selector = '[name=accountCode' + this.data.suffix + ']';
    $(selector).typeahead({
      highlight: true
    }, {
      displayKey: 'code',
      source: this.datasource.ttAdapter(),
      templates: {
        suggestion: function(item) {
          return item.displayName;
        }
      }
    }).on('typeahead:selected', function(evt, object) {
      this.set({
        account: object.code,
        label: object.label
      });
    }.bind(this));
  };

  var counter = 0;
  var AccountPicker = Ractive.extend({
    template: '<div class="account-picker">' +
    '<input id="accountCode{{suffix}}" name="accountCode{{suffix}}" type="text" placeholder="Identifikation" value="{{account}}"/>' +
    '<input id="accountLabel{{suffix}}" name="accountLabel{{suffix}}" type="text" value="{{label}}">' +
    '</div>',
    isolated: true,
    complete: function() {
      this.set({
        suffix: this.counter
      });
      setupAutocompletion.bind(this).call();
    },
    init: function() {
      this.counter = counter + 1;
      counter += 1;
    },
    data: {
      account: '',
      label: ''
    }
  });

  Ractive.components['account-picker'] = AccountPicker;
});
