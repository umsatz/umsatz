define([], function() {
  'use strict';

  return {
    totalIncome: function() {
      var total = 0;
      this.filter(function(position) {
        return position.isIncome();
      }).forEach(function(position) {
        total += position.signedTotalAmountCents();
      });
      return total;
    },
    totalExpense: function() {
      var total = 0;
      this.filter(function(position) {
        return !position.isIncome();
      }).forEach(function(position) {
        total += position.signedTotalAmountCents();
      });
      return total;
    },
    totalAmount: function() {
      var total = 0;
      this.forEach(function(position) {
        total += position.signedTotalAmountCents();
      });
      return total;
    },
    totalIncomeVatAmount: function() {
      var incomeVat = 0;
      this.filter(function(position) {
        return position.isIncome();
      }).forEach(function(position) {
        incomeVat += position.totalVatAmountCents();
      });
      return incomeVat;
    },
    totalExpenseVatAmount: function() {
      var expenseVat = 0;
      this.filter(function(position) {
        return !position.isIncome();
      }).forEach(function(position) {
        expenseVat += position.totalVatAmountCents();
      });
      return expenseVat;
    },
    totalVatAmount: function() {
      var totalVat = 0;
      this.forEach(function(position) {
        if (position.isIncome()) {
          totalVat += position.totalVatAmountCents();
        } else {
          totalVat -= position.totalVatAmountCents();
        }
      });
      return totalVat;
    }
  };
});
