// detect.ac advanced theme - communicates with the trusted shell via postMessage.
// The shell validates every action and routes it to native bindings (submit PIN,
// agree/decline TOS, exit, drag the window). The sandbox prevents any other I/O.

var bridge = {
  send: function (op, data) {
    parent.postMessage(Object.assign({ kind: 'detect.action', op: op }, data || {}), '*');
  },
  agreeTos:   function ()   { this.send('tos.agree'); },
  declineTos: function ()   { this.send('tos.decline'); },
  submitPin:  function (v)  { this.send('pin.submit', { pin: String(v || '').replace(/\D/g, '').slice(0, 8) }); },
  exit:       function ()   { this.send('app.exit'); },
  drag:       function ()   { this.send('app.drag'); }
};

// State updates pushed from the shell: { kind: 'detect.state', stage, progress, status, pinError }
window.addEventListener('message', function (ev) {
  if (!ev.data || ev.data.kind !== 'detect.state') return;
  if (ev.data.stage) {
    document.querySelectorAll('.stage').forEach(function (s) {
      s.classList.toggle('active', s.dataset.stage === ev.data.stage);
    });
  }
  if (typeof ev.data.status === 'string') {
    var el = document.getElementById('scanner-status'); if (el) el.textContent = ev.data.status;
  }
});

// TOS scroll → enable agree button when fully read
var tosBody  = document.getElementById('tos-body');
var btnAgree = document.getElementById('btn-agree');
var tosHint  = document.getElementById('tos-hint');
var tosPill  = document.getElementById('tos-progress-pill');

if (tosBody) tosBody.addEventListener('scroll', function () {
  var pct = tosBody.scrollTop / (tosBody.scrollHeight - tosBody.clientHeight);
  var clamped = Math.min(100, Math.round(pct * 100));
  if (tosPill) tosPill.textContent = clamped + '%';
  if (clamped >= 100) {
    if (btnAgree) btnAgree.disabled = false;
    if (tosHint)  tosHint.style.visibility = 'hidden';
  }
});

// TOS buttons
if (btnAgree) btnAgree.addEventListener('click', function () { bridge.agreeTos(); });

var btnDecline = document.getElementById('btn-decline');
if (btnDecline) btnDecline.addEventListener('click', function () { bridge.declineTos(); });

// Close / exit buttons
['btn-close', 'btn-err-exit', 'btn-self-exit'].forEach(function (id) {
  var el = document.getElementById(id);
  if (el) el.addEventListener('click', function () { bridge.exit(); });
});

// Topbar drag - ask the shell to start the native Win32 caption move loop.
var bar = document.querySelector('.topbar');
if (bar) bar.addEventListener('mousedown', function (e) {
  if (e.button !== 0) return;
  if (e.target.closest('button')) return;
  bridge.drag();
});
