function ShowData() {
    let showName = document.getElementById('input-name').value;
    let showEmail = document.getElementById('input-email').value;
    let showPhone = document.getElementById('input-number').value;
    let showSubject = document.getElementById('input-subject').value;
    let showMessage = document.getElementById('input-message').value;


    if (showName =='') {
        return alert('nama harus diisi')
    }
    if (showEmail == '') {
        return alert('email harus diisi')
    }
    if (showPhone == '') {
        return alert('no.Telepon harus diisi')
    }
    if (showSubject == '') {
        return alert('Subject harus diisi')
    }

    console.log(showEmail)
    console.log(showName)
    console.log(showSubject)

    let emailReceiver = 'rezkirahman0509@gmail.com'

    let a = document.createElement('a');
    a.href = `mailto: ${emailReceiver}?subject: ${showSubject}&body= hello, my name is ${showName}, ${showMessage}`

    a.click()
}

