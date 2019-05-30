const bcrypt = require('bcryptjs');
let password = "123456";
const sha256 = require('js-sha256').sha256;
password=sha256(password);
console.log(password);
let hash =  "$2a$10$la75mLwUDCkxwdMNdOBaS.UHdjo3MD2iESfAmNTM1/h2vgHkFTdYm";
let result = bcrypt.compareSync(password,hash);
console.log(result);