import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import MFab from '@material-ui/core/Fab';
import Dialog from '@material-ui/core/Dialog';
import Fade from '@material-ui/core/Fade';
import AddIcon from '@material-ui/icons/Add';
import Backdrop from '@material-ui/core/Backdrop';

const useStyles = makeStyles((theme) => ({
    modal: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
}));

interface IProps {
    modal: React.ReactNode;
    includeFade?: boolean;
}

function buildFade(component: React.ReactNode, fadeProps: any, includeFade?: boolean) {
    if (!includeFade) {
        return component;
    }

    return (
        <Fade {...fadeProps} >{component}</Fade>
    );

}

export function Fab(props: IProps) {
    const Modal = props.modal;
    const classes = useStyles();
    const [open, setOpen] = useState(false);
    const handleModalOpen = () => {
        setOpen(true);
    };

    const handleModalClose = (data: any) => {
        setOpen(false);
    };
    return (
        <div>
            <MFab
                color="primary"
                style={{ position: 'fixed', right: '15px', bottom: '15px' }}
                onClick={handleModalOpen}
            >
                <AddIcon />
            </MFab>
            <Dialog
                aria-labelledby="transition-modal-title"
                aria-describedby="transition-modal-description"
                className={classes.modal}
                open={open}
                onClose={handleModalClose}
                closeAfterTransition
                BackdropComponent={Backdrop}
                BackdropProps={{
                    timeout: 500,
                }}
            >
                {buildFade(Modal, { in: open}, props.includeFade )}
            </Dialog>
        </div>
    )
}