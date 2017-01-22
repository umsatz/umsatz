'use strict';

define([
  'intern/chai!expect',
  'intern!bdd',

  'intern/order!scripts/models/account',
], function(expect, bdd, Account) {
  var account = null;
  bdd.describe('Account', function() {

    bdd.beforeEach(function() {
      account = new Account({
        label: 'Sample',
        code: '200'
      });
    });

    bdd.describe('#url', function() {

      bdd.it('adds trailing / for new accounts', function() {
        expect(account.url()).to.eql('/api/accounts/');
      });
      bdd.it('adds /id for existing accounts', function() {
        account.set('id', 42);
        expect(account.url()).to.eql('/api/accounts/42');
      });

    });

    bdd.describe('#toJSON', function() {

      bdd.it('contains a displayName', function() {
        var data = account.toJSON();
        expect(data.displayName).to.eql('Sample <200>');
      });

    });
  });
});