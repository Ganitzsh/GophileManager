var socket
var host = 'http://localhost:9000'

$(document)
	.ready(function () {
		socket = io(host);
		socket.on('notif action done', function (v) {
			$('#' + v.alert)
				.remove()
			newNotification('done', v.message)
			if (v.reload) {
				fetchContent('current')
			}
		})
		socket.on('notif action error', function (v) {
			$('#' + v.alert)
				.remove()
			newNotification('error', v.message)
			if (v.reload) {
				fetchContent('current')
			}
		})
		socket.on('notif action progress', function (v) {
			console.log(v);
			$('#' + v.alert)
				.remove()
			$('#notifications')
				.append(generateProgressAlert(v.message, v.progress, v.alert))
			if (v.reload) {
				fetchContent('current')
			}
		})
		$('.sk-fading-circle')
			.hide()
		fetchContent('current')
	});

function generateProgressAlert(text, progress, id) {
	return `
	<div class="alert alert-info alert-dismissible" role="alert" id="${id}">
		<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
		${text}
		<div class="progress">
				<div class="progress-bar progress-bar-striped active" role="progressbar" style="width: ${progress}%">
				</div>
		</div>
	</div>
	`
}

function generateAlert(type, text, id) {
	return `
	<div class="alert alert-${type} alert-dismissible" role="alert" id="${id}">
		<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
		${text}
	</div>
	`
}

function newNotification(type, message, arg) {
	var alert

	var uuid = UUID.generate()
	var domID = 'alert-' + uuid
	switch (type) {
	case 'progress':
		if (arg === undefined) {
			alert = generateAlert('danger', 'Invalid progress alert!', domID)
		} else {
			alert = generateProgressAlert(message, arg.progress, domID)
		}
		break;
	case 'done':
		alert = generateAlert('success', message, domID)
		break;
	case 'error':
		alert = generateAlert('danger', message, domID)
		break;
	}
	$('#notifications')
		.append(alert)
	console.log(domID);
	return domID
}

function loading() {
	$('#files')
		.hide()
	$('.sk-fading-circle')
		.show();
}

function doneLoading() {
	$('#files')
		.show()
	$('.sk-fading-circle')
		.hide();
}

function reloadFiles(newFiles) {
	$('html, body')
		.animate({
			scrollTop: 0
		}, 'fast');
	$('#files')
		.html(newFiles)
	$('[data-toggle="tooltip"]')
		.tooltip()
}

function fetchContent(where) {
	var str = '/app/files/' + where
	var uri = encodeURI(str)
	loading()
	fetch(uri, {
			credentials: 'same-origin' // Add cookie
		})
		.then(function (response) {
			if (response.ok) {
				response.text()
					.then(function (html) {
						reloadFiles(html)
						doneLoading()
					});
			} else {}
		})
		.catch(function (error) {});
}

function watch(name) {
	window.location = "/app/watch/" + name;
}

function convertToMP4(name) {
	bootbox.confirm("<h3>Convert</h3><p><strong>" + name + "</strong></p><p>Convert it to MP4?</p>", function (result) {
		if (result) {
			var alertID = newNotification('progress', 'Converting <strong>' + name + '</strong>', {
				progress: 100
			})
			var str = '/app/convert/' + name
			str += '?alert_id=' + alertID
			var uri = encodeURI(str)
			fetch(uri, {
					credentials: 'same-origin'
				})
				.then(function (response) {
					if (response.ok) {} else {}
				})
				.catch(function (error) {});
		}
	})
}

function trashFile(name) {
	bootbox.confirm("<h3>Warning!</h3><p>Move <strong>" + name + "</strong> to trash?</p>", function (result) {
		if (result) {
			var alertID = newNotification('progress', 'Moving <strong>' + name + '</strong> to trash', {
				progress: 100
			})
			var str = '/app/trash/file/' + name
			str += '?alert_id=' + alertID
			var uri = encodeURI(str)
			fetch(uri, {
					credentials: 'same-origin'
				})
				.then(function (response) {
					if (response.ok) {} else {}
				})
				.catch(function (error) {});
		}
	})
}

function emptyTrash() {
	bootbox.confirm("<h3>Warning!</h3><p>Empty the trash?</p><strong>All data will be definitely lost</strong>", function (result) {
		if (result) {
			var alertID = newNotification('progress', 'Moving <strong>' + name + '</strong> to trash', {
				progress: 100
			})
			var str = '/app/trash/empty'
			str += '?alert_id=' + alertID
			var uri = encodeURI(str)
			fetch(uri, {
					credentials: 'same-origin'
				})
				.then(function (response) {
					if (response.ok) {} else {}
				})
				.catch(function (error) {});
		}
	})
}

function previewImage(name) {
	$('#imgModalSrc')
		.attr('src', '/app/file/' + name)
	$('#imageModal')
		.modal('show')
}

function deleteFile(name) {
	bootbox.confirm("<h3>Warning!</h3><p>Delete <strong>" + name + "</strong>?</p>", function (result) {
		if (result) {
			var alertID = newNotification('progress', 'Deleting <strong>' + name + '</strong>', {
				progress: 100
			})
			var str = '/app/delete/' + name
			str += '?alert_id=' + alertID
			var uri = encodeURI(str)
			fetch(uri, {
					credentials: 'same-origin'
				})
				.then(function (response) {
					if (response.ok) {} else {}
				})
				.catch(function (error) {});
		}
	})
}

function downloadFile(name) {
	var str = '/app/download/' + name
	var uri = encodeURI(str)
	console.log(uri);
	fetch(uri, {
			credentials: 'same-origin'
		})
		.then(function (response) {
			if (response.ok) {
				var mime = response.headers.get('Content-Type')
				response.blob()
					.then(function (blob) {
						download(blob, name, mime)
					});
			} else {}
		})
		.catch(function (error) {
			console.log('There has been a problem with your fetch operation: ' + error.message);
		});
}

function resetCompressionModal() {
	$('#fileName')
		.text('')
	$('#archiveFileName')
		.val('')
	$('#compressionModal')
		.modal('hide')
}

function setToCompress(name) {
	$('#fileName')
		.text(name)
	$('#compressionModal')
		.modal('show')
}

function validateCompressionForm() {
	var name = $('#archiveFileName')
		.val()
	console.log(name);
	if (!name) {
		$('#compressionFormGroup')
			.addClass('has-error')
		$('#nameHelp')
			.removeClass('hidden')
	} else {
		var file = $('#fileName')
			.text()
		compressAndDownload(name, file)
	}
}

function compressAndDownload(name, target) {
	var alertID = newNotification('progress', 'Compressing <strong>' + target + '</strong>', target)
	resetCompressionModal()
	var str = '/app/compress/' + target + '/name/' + name
	str += '?alert_id=' + alertID
	var uri = encodeURI(str)
	fetch(uri, {
			credentials: 'same-origin'
		})
		.then(function (response) {
			if (response.ok) {} else {}
		})
		.catch(function (error) {});
}
