import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardActions from '@material-ui/core/CardActions';
import TextField from '@material-ui/core/TextField';
import CircularProgress from '@material-ui/core/CircularProgress';
import Popover from '@material-ui/core/Popover';
import { Project, ProjectAPI } from "../../../services/project";
import { ColorPicker } from '../../../components/color-picker';

const useStyles = makeStyles((theme) =>({
    root: {
        minHeight: 400,
        minWidth: 400
    },
    title: {
        top: 1,
    },
    actionButtons: {
        top: '145px',
        position: 'relative'
    },
    progress: {
        display: 'flex',
        '& > * + *': {
          marginLeft: theme.spacing(2),
        },
    },
  }),
);

interface IProps {
    Callback(project: Project): void;
    ProjectAPI: ProjectAPI;
}


export function ProjectModal(props: IProps) {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [showProgress, setShowProgress] = useState(false);
    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const [projectColor, setProjectColor] = useState("");

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const open = Boolean(anchorEl);
    const onSubmit = function (event: React.FormEvent) {
        event.preventDefault();
    }

    const onSave = async () => {
        const proj = {
            metadata: {
                name,
                description,
                id: '',
                color: projectColor,
                imageUrl: ''
            },
            tasks: [],
        }
        const project = await props.ProjectAPI.create(proj)
        setShowProgress(true)
        setTimeout(() => {
            props.Callback(project)
        }, 1000)
    }

    const onColorSet = (hex: string) => {
        setProjectColor(hex);
        handleClose();
    }

    const classes = useStyles();
    if (showProgress) {
        return( 
            <div className={classes.progress}>
                <CircularProgress></CircularProgress>
            </div>
        )
    }
    return (
        <form noValidate autoComplete="off" onSubmit={onSubmit}>
            <Card className={classes.root}>
                <TextField value={name} onChange={(ev: any) => setName(ev.target.value)} fullWidth label="Name" variant="filled" />
                <CardContent>
                    <TextField
                        label="Description"
                        multiline
                        value={description}
                        onChange={(ev: any) => setDescription(ev.target.value)}
                        rows={4}
                        fullWidth
                        variant="filled"
                    />
                </CardContent>
                <CardActions className={classes.actionButtons}>
                    <Button type="submit" size="small" color="primary" onClick={onSave}>
                        Save
                    </Button>
                    <div>
                        <Button variant="contained" color="primary" onClick={handleClick}>Pick Color</Button>
                        <Popover
                            id={Date.now().toString()}
                            open={open}
                            anchorEl={anchorEl}
                            onClose={handleClose}
                            anchorOrigin={{
                                vertical: 'bottom',
                                horizontal: 'center',
                            }}
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'center',
                            }}
                        >
                            <ColorPicker onColorSet={onColorSet}/>
                        </Popover>
                    </div>
                </CardActions>
            </Card>
        </form>
    )
}