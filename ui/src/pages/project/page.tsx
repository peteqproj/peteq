import React, { useState , useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Chip from '@material-ui/core/Chip';
import Typography from '@material-ui/core/Typography';
import { RouteComponentProps } from "react-router-dom";
import { ProjectAPI } from './../../services/project'
import { ProjectViewAPI, ProjectView } from './../../services/view.project'

const useStyles = makeStyles({
  root: {
    display: 'flex'
  },
  media: {
    height: 320,
    width: 320
  },
});


interface IProps extends RouteComponentProps {
  id?: string;
  ProjectAPI: ProjectAPI;
  ProjectViewAPI: ProjectViewAPI;
}

export function ProjectPage (props: IProps) {
  const classes = useStyles();
  const [state, setState] = useState({
    metadata: {
      name: '',
      id: '',
      description: '',
      color: '',
      imageUrl: '',
    }, tasks: []
  } as ProjectView);
  const id = (props.match.params as any)['id'];
  useEffect(() => {    
    const fetch = async () => {
      const prj = await props.ProjectViewAPI.get(id)
      setState(prj)
    }
    fetch();
  
  }, [props.ProjectAPI, id]);
  return (
    <Card className={classes.root}>
      <CardActionArea>
        <CardContent>
          <CardHeader component={() => (<Chip style={{ width: '70px', backgroundColor: state.metadata.color }}/>)}/>
          <Typography gutterBottom variant="h3" component="h2">
            {state.metadata.name}
          </Typography>
          <Typography variant="body1" color="textSecondary" component="p">
            {state.metadata.description}
          </Typography>
          <Typography gutterBottom variant="h6" component="h2">
            Tasks:
          </Typography>
          <Typography variant="body2" color="textSecondary" component="p">
            {state.tasks.map((t) => (<div>{t.metadata.name}</div>))}
          </Typography>
        </CardContent>
      </CardActionArea>
      {state.metadata.imageUrl !== '' && <CardMedia
          className={classes.media}
          image={state.metadata.imageUrl}
          title="Contemplative Reptile"
        />}
    </Card>
  );
}

