var _elqQ = _elqQ || [];
_elqQ.push(['elqSetSiteId', '88570519']);
_elqQ.push(['elqTrackPageView']);

(function () {
  function async_load() {
    var s = document.createElement('script'); s.type = 'text/javascript'; s.async = true;
    s.src = '//img04.en25.com/i/elqCfg.min.js';
    var x = document.getElementsByTagName('script')[0]; x.parentNode.insertBefore(s, x);
  }
  if (window.addEventListener) window.addEventListener('DOMContentLoaded', async_load, false);
  else if (window.attachEvent) window.attachEvent('onload', async_load);
})();

var timerId = null, timeout = 5;
function WaitUntilCustomerGUIDIsRetrieved() {
  if (!!(timerId)) {
    if (timeout == 0) {
      return;
    }
    if (typeof this.GetElqCustomerGUID === 'function') {
      var form = document.forms["ContactUs_TemplateForm"];
      if (form) form.elements["elqCustomerGUID"].value = GetElqCustomerGUID();
      return;
    }
    timeout -= 1;
  }
  timerId = setTimeout("WaitUntilCustomerGUIDIsRetrieved()", 500);
  return;
}
window.onload = WaitUntilCustomerGUIDIsRetrieved;
_elqQ.push(['elqGetCustomerGUID']);