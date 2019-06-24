//SPDX-License-Identifier: Apache-2.0

var account = require('./controller.js');

module.exports = function(app){

  app.get('/get_user/:id', function(req, res){
    account.get_user(req, res);
  });
  app.get('/add_user/:account', function(req, res){
    account.add_user(req, res);
  });
  app.get('/trans_money/:trans', function(req, res){
    account.trans_money(req, res);
  });
}
