define([
  'backbone',
  'models/account',
], function(Backbone, Account) {
  'use strict';

  return Backbone.Collection.extend({
    model: Account,
    url: '/api/accounts/',

    /**
     * creates or updates an account, identified by its code.
     */
    upsert: function(attrs) {
      var account;
      if ((account = this.findWhere({
          code: attrs.code
        })) !== undefined) {
        account.set(attrs);
        return account.save();
      } else {
        return this.create(attrs);
      }
    }
  });
});
