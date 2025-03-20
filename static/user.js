async function profile(){
    const token = localStorage.getItem('token')
    try {
        const response = await fetch(account_url + "/profile", {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        const data = await response.json();
        console.log(data);
        if (response.ok) {
            return data;
        }else{
            directToLogin();
        }
    } catch (error) {
        console.error('ユーザー情報の取得に失敗しました:', error);
        directToLogin();
    }
}
function logout() {
    localStorage.removeItem('token');
    directToLogin();
}

async function login(username, password) {
    const response = await fetch(account_url + "/login", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username: username, password: password })
    });

    const data = await response.json();
    if (response.ok){
        console.log(data);
        localStorage.setItem('token', data.token);
        //alert("ログイン成功！トップページへ移動します。");
        return {"ok": response.ok,"data":data};
    }else{
        console.error(data.error);
        return {"ok": response.ok,"error":data.error};
    }
}

async function register(username, password) {
    var error_message;
    const response = await fetch(account_url + "/register", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            password: password,
        }),
    });

    const data = await response.json();
    //console.log(data);
    if (response.ok){
        await login(username, password);
        return {"ok":true, "data":data};
    }else{
        error_message = data.error;
        console.error(error_message);
        return {"ok":false, "error":error_message};
    }
}

function directToLogin(){
    window.location.href = '/login';
}