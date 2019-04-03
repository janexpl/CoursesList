jQuery.extend (jQuery.fn.dataTableExt.oSort, {
  'cert-number-pre' : function( a ) {
    if (!a) { return 0 }
      var i, item;
      var m = a.split("/");
      //pobranie numeru certyfikatu
      item = m[0];
      console.log(item);
      return parseInt(item);
  },
  'cert-number-asc': function ( a, b ) {
    console.log(a, b);
    return (( a < b ) ? -1 : ((a > b) ? 1 : 0));
  },
  'cert-number-desc': function ( a, b ) {
    console.log(a, b);
    return (( a < b ) ? 1 : ((a > b) ? -1 : 0));
  }  
});
