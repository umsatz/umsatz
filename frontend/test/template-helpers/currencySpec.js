'use strict';

define([
  'intern/chai!expect',
  'intern!bdd',

  'intern/order!scripts/template-helpers/currency',
], function(expect, bdd, currency) {
  /* jshint -W030 */

  bdd.describe('currency', function() {
    bdd.it('defaults currency to EUR', function() {
      expect(currency(200)).to.eql('2,00 €');
    });

    bdd.it('formattes USD as $', function() {
      expect(currency(200, {
        currency: 'USD'
      })).to.eql('2,00 $');
    });

    bdd.it('formattes EUR as €', function() {
      expect(currency(200, {
        currency: 'EUR'
      })).to.eql('2,00 €');
    });

    bdd.it('supports custom cent metrics', function() {
      expect(currency(20, {
        cents: 10
      })).to.eql('2,00 €');
    });
  });
});