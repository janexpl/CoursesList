$(document).ready(function() {
  $("#coursesTable").DataTable();
  $("#studentsTable").DataTable();
  $("#compTable").DataTable();
  $("#certsTable").DataTable();
  $(".dataTables_length").addClass("bs-select");
});

$("[data-toggle=confirmation]").confirmation({
  rootSelector: "[data-toggle=confirmation]"
  // other options
});

$(function() {
  $("#flash")
    .delay(500)
    .fadeIn("normal", function() {
      $(this)
        .delay(2500)
        .fadeOut();
    });
});

jQuery(function($) {
  $("form[data-async]").on("submit", function(event) {
    var $form = $(this);
    var $target = $($form.attr("data-target"));
    $.ajax({
      type: $form.attr("method"),
      url: $form.attr("action"),
      data: $form.serialize(),

      success: function(data, status) {
        $target.html(data);
      }
    });

    event.preventDefault();
  });
});
