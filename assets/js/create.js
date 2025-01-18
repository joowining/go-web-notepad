// password duplicate check 

const firstPwd = document.getElementById("first-pwd")
const secondPwd = document.getElementById("second-pwd")
const pwdCheckBtn = document.getElementById("pwd-check")

pwdCheckBtn.addEventListener("click",()=>{
    console.log("first", firstPwd.value);
    console.log("second", secondPwd.value);
    if (firstPwd.value === ""){
        alert("password is empty")
        return 
    }
    if (firstPwd.value === secondPwd.value){
        alert("password matched! go next");
    }else {
        alert("password mismatched! try again");
    }
})


// 사용자의 입력 확인 
const userId = document.getElementById("user-id");
const userName = document.getElementById("user-name");
const userEmail = document.getElementById("user-email");
const submitUserInfoBtn = document.querySelector(".sign-up-btn");

function isValidEmail(email) {
    // 이메일 형식을 확인하는 정규표현식
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailRegex.test(email);
}

function checkEmptyString(value) {
    if (value === ""){
        return false;
    }
    return true;
}

submitUserInfoBtn.addEventListener("click", (event) => {
   if( !checkEmptyString(userId.value)){
        event.preventDefault();
        alert("input your id correctly");
        return;
   }else if(!checkEmptyString(userName.value)){
        event.preventDefault();
        alert("input your name correctly");
        return; 
   }else if(!checkEmptyString(firstPwd.value)){
        event.preventDefault();
        alert("input your password correctly");
        return;
   }else if(!isValidEmail(userEmail.value)){ 
        event.preventDefault();
        alert("input your email as email form ");
        return;
   }else {
        return;
   }
})

