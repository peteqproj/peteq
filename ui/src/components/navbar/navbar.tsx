import React from 'react'
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';


interface IProps {
    className: string;
}

export function Navbar(props: IProps) {
    return (
        <AppBar position="fixed" className={props.className}>
            <Toolbar>
                <Typography variant="h6" noWrap></Typography>
            </Toolbar>
        </AppBar>
    );
}