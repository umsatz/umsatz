'use strict';

define([
  'intern/chai!expect',
  'intern!bdd',

  'intern/order!scripts/models/fiscalPeriod',
], function(expect, bdd, FiscalPeriod) {
  var fiscalPeriod = null;
  bdd.describe('FiscalPeriod', function() {

    bdd.beforeEach(function() {
      fiscalPeriod = new FiscalPeriod({
        label: 'Sample',
        code: '200'
      });
    });

    bdd.describe('#url', function() {

      bdd.it('adds trailing / for new fiscalPeriod', function() {
        expect(fiscalPeriod.url()).to.eql('/api/fiscalPeriods/');
      });
      bdd.it('adds /id for existing fiscalPeriods', function() {
        fiscalPeriod.set('id', 42);
        expect(fiscalPeriod.url()).to.eql('/api/fiscalPeriods/42');
      });

    });
  });
});