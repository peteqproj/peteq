import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Paper from '@material-ui/core/Paper';
import { UserAPI } from './../../services/user';


const useStyles = makeStyles((theme) => ({
    root: {
        backgroundColor: 'gray',
        flexGrow: 1,
    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
        height: '400px',
    },
    formItem: {
        textAlign: 'center',
        padding: theme.spacing(2),
        width: '100%'
    }
}));

interface IProps {
    UserAPI: UserAPI;
    onLogin(): void;
}

export function LoginPage(props: IProps) {
    const classes = useStyles()
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const doLogin = async () => {
        const res = await props.UserAPI.login(email, password);
        localStorage.setItem("API_TOKEN", res.token)
        props.UserAPI.updateDefaultsWithAPIToken(res.token)
        props.onLogin()
    }
    return (
        <Grid
            className={classes.root}
            container spacing={3}
            direction="row"
            justify="space-between"
            alignItems="stretch">
            <Grid item xs={12} style={{height: '200px'}} />
            <Grid item xs={3} />
            <Grid item xs={6}>
                <Paper className={classes.paper}>
                    <TextField className={classes.formItem} id="email" label="Email" onChange={ev => setEmail(ev.target.value)} />
                    <TextField type={"password"} className={classes.formItem} id="password" label="Password" onChange={ev => setPassword(ev.target.value)} />
                    <Button className={classes.formItem} variant="contained" color="primary" onClick={doLogin}>
                        Login
                    </Button>
                </Paper>
            </Grid>
            <Grid item xs={3} />
            <Grid item xs={12} />
        </Grid>
    )
}

export function getAPIToken(): string {
    const token = localStorage.getItem('API_TOKEN')
    if (!token || token === "") {
        return "";
    }
    return token
}