
// 입력한 dom요소에 추가할 dom 요소를 어떤 내용과 함께 추가
function addDomEl(parent,  elName, content) {    
    let parentEl = document.querySelector(parent);
    let childEl = document.createElement(elName);

    childEl.content = content;
    parentEl.child = childEl;
}