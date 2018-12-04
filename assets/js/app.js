$(document).ready(function () {
    $('#coursesTable').DataTable();   
    $('#studentsTable').DataTable();
    $('#compTable').DataTable();
    $('.dataTables_length').addClass('bs-select');
    
    });

    
$('[data-toggle=confirmation]').confirmation({
    rootSelector: '[data-toggle=confirmation]',
    // other options
});



jQuery(function($) {
  $('form[data-async]').on('submit', function(event) {
      var $form = $(this);
      var $target = $($form.attr('data-target'));
      // var nipsum = $form.attr('nip');
      // parseInt(nip.charAt(0))*6+parseInt(nip.charAt(1))*5+parseInt(nip.charAt(2))*7+parseInt(nip.charAt(3))*2+parseInt(nip.charAt(4))*3+parseInt(nip.charAt(5))*4+parseInt(nip.charAt(6))*5+parseInt(nip.charAt(7))*6+parseInt(nip.charAt(8))*7;
			// reszta = nipsum % 11;
			// if(reszta != parseInt(nip.charAt(9))){
			// 	console.log(reszta)
			// 	alert("NIP niepoprawny");
      //   return false;
			// }
      $.ajax({
          type: $form.attr('method'),
          url: $form.attr('action'),
          data: $form.serialize(),

          success: function(data, status) {
              $target.html(data);
          }
      });

      event.preventDefault();
  });
});



  

