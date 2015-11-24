var activeTube = '';

$(function() {
    $.get("/host", function(data) {
        $("title").append(data);
    });
    updateStats();
    setInterval(updateStats, 1000);

    $("#stats").on('click', 'button.peek-ready', function(e) {
        var tube = $(e.target).data('tube');
        var job = $.get("/peek-ready/" + tube, function(data) {
            showPeekModal(data);
        });
    });

    $("#stats").on('click', 'button.peek-buried', function(e) {
        var tube = $(e.target).data('tube');
        var job = $.get("/peek-buried/" + tube, function(data) {
            showPeekModal(data);
        });
    });

    $("#stats").on('click', 'button.edit', function(e) {
        var tube = $(e.target).data('tube');
        $("#tube-modal").data('tube', tube);
        $("#tube-modal").modal();
    });

    $("#kick-all").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/kick-all/" + tube, function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#kick").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/kick/" + tube + "/" + $("#kick-count").val(), function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#bury-all").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/bury-all/" + tube, function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#bury").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/bury/" + tube + "/" + $("#bury-count").val(), function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#drain-all-buried").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/drain-all-buried/" + tube, function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#drain-buried").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/drain-buried/" + tube + "/" + $("#drain-buried-count").val(), function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#drain-all-ready").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/drain-all-ready/" + tube, function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#drain-ready").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/drain-ready/" + tube + "/" + $("#drain-ready-count").val(), function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#drain-all-buried").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/drain-all-buried/" + tube, function() {
            $("#tube-modal").modal('hide');
        });
    });

    $("#insert").on('click', function() {
        var tube = $("#tube-modal").data('tube');
        $.post("/insert/" + tube, {
            data: $("#insert-data").val().trim()
        }, function() {
            $("#tube-modal").modal('hide');
        });
    });
});


function updateStats() {
    $.get("/stats", function(data) {

        var currentTubes = [];
        $.each(data, function(tube, stats) {
            currentTubes.push(tube);
            var rowSelector = "#" + tube + "-stats";
            if ($(rowSelector).length == 0) {
                $("#stats").append("<tr id=\"" + tube + "-stats\"></tr>");
            }
            var rowHtml = "<td>" + tube + "</td>";
            rowHtml += "<td>" + stats["current-jobs-buried"] + "</td>";
            rowHtml += "<td>" + stats["current-jobs-ready"] + "</td>";
            rowHtml += "<td>" + stats["current-jobs-reserved"] + "</td>";
            rowHtml += "<td>" + stats["current-waiting"] + "</td>";
            rowHtml += "<td>" + stats["total-jobs"] + "</td>";
            rowHtml += "<td>";
            rowHtml += "<button data-tube=\"" + tube + "\" class=\"edit btn btn-info btn-mini\"><i class=\"icon-white icon-pencil\"></i></button> ";
            rowHtml += "<button data-tube=\"" + tube + "\" class=\"peek-ready btn btn-primary btn-mini\"><i class=\"icon-white icon-eye-open\"></i></button> ";
            rowHtml += "<button data-tube=\"" + tube + "\" class=\"peek-buried btn btn-primary btn-mini\"><i class=\"icon-white icon-eye-close\"></i></button> ";
            rowHtml += "</td>";
            $(rowSelector).html(rowHtml);
        });


        $("#stats").find("tr").each(function(idx, el) {
            var id = $(el).attr('id');
            if (id != "stats-header") {
                var tubeName = id.replace("-stats", "");
                for (var i = 0; i < currentTubes.length; i++) {
                    if (currentTubes[i] == tubeName) {
                        return
                    }
                }
                $(el).remove();
            }
        });
    });
}

function showPeekModal(data) {
    $("#job-id").val(data.job_id);
    $("#job-data").html(data.job_data);
    $("#peek-modal").modal();
}