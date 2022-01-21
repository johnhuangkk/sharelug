$(function () {
  let prevScrollpos = window.pageYOffset;
  window.onscroll = function () {
    let currentScrollPos = window.pageYOffset;
    if (currentScrollPos < 0) {
      currentScrollPos = 0;
    }
    if (prevScrollpos >= currentScrollPos) {
      document.getElementById("pageHeader").style.top = "0";
    } else {
      document.getElementById("pageHeader").style.top = "-50px";
    }
    prevScrollpos = currentScrollPos;
  }
});
