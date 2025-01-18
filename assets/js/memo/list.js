const URL = "http://localhost:8000"

// 1. Username 요청하고 가져오기 
fetch(URL + '/user')
    .then(response => response.json())
    .then(data => {
        console.log(data);
        const userHeader = document.querySelector('.user-name');
        userHeader.textContent = data.user;
    }).catch(err => console.log(err))
// 2. with all delete options

const deleteOptionBtn = document.querySelector('.deleteOptionBtn')
deleteOptionBtn.addEventListener('click', changeDeleteState)
let deleteOption = false;

function changeDeleteState(){
    deleteOption = !deleteOption;

    changeDeleteCheckBox(deleteOption);
    if (deleteOption){
        deleteOptionBtn.textContent = "취소"        
    }else {
        deleteOptionBtn.textContent = "일괄 삭제"
    }

}

function changeDeleteCheckBox( option ){
    const elList = document.querySelector('.block-memo-list');

    if( option ){
        for(let i = 0; i < elList.children.length; i++){
            const deleteBox = elList.children[i].children[1].children[3];
            deleteBox.setAttribute('style','display:inline-block')
        }
    }else {
        for(let i = 0; i < elList.children.length; i++){
            const deleteBox = elList.children[i].children[1].children[3];
            deleteBox.setAttribute('style','display:none')
        }
    }
} 

// 3. get every 10 memos from server and create its list  
fetch(URL + '/memo/listitem')
    .then(response => response.json())
    .then(data => {
        // 데이터의 내용들은 
        // 1.타이틀 
        // 2.memoId
        // 데이터의 갯수는 10개 
        for( let i = 0; i < data.memos.length; i++){
            createMemoListItem(data.memos[i]);
        }
    })
    .catch(error=>{
        console.log(error);
    })
const listContent = document.querySelector('.block-memo-list');

function createMemoListItem( data){
    let listDiv = document.createElement('div');
    listDiv.setAttribute('class','list-memo-item');
    let title = document.createElement('p');
    title.setAttribute('class','title-text');
    title.innerText = data.title;

    let btnDiv = document.createElement('div');
    let openBtn = document.createElement('button');
    openBtn.innerText = '열기';
    openBtn.setAttribute('class','button ' + data.id);
    let editBtn = document.createElement('button');
    editBtn.innerText = '수정';
    editBtn.setAttribute('class','button ' + data.id);
    let deleteBtn = document.createElement('button');
    deleteBtn.innerText = '삭제';
    deleteBtn.setAttribute('class','button ' + data.id);

    listDiv.appendChild(title);
    btnDiv.appendChild(openBtn);
    btnDiv.appendChild(editBtn);
    btnDiv.appendChild(deleteBtn);
    listDiv.appendChild(btnDiv);
    
    deleteCheckBox = document.createElement('input')
    deleteCheckBox.setAttribute('class','delete-box');
    deleteCheckBox.setAttribute('type','checkbox');
    deleteCheckBox.setAttribute('name','delete-box');
    // deleteCheckBox.setAttribute('value', memoId);
    btnDiv.appendChild(deleteCheckBox); 

    listContent.appendChild(listDiv);
}

const createElBtn = document.querySelector('.createEl');
createElBtn.addEventListener('click', createMemoListItem)

// 4. pagination 