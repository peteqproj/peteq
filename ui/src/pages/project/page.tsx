import React, { useState, useEffect } from 'react'
import { createStyles, makeStyles, useTheme, Theme } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import Input from '@material-ui/core/Input';
import MenuItem from '@material-ui/core/MenuItem';
import { ProjectAPI, Project } from './../../services/project'
import { RouteComponentProps } from "react-router-dom";


const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    formControl: {
      margin: theme.spacing(1),
      minWidth: 120,
      maxWidth: 300,
    },
    chips: {
      display: 'flex',
      flexWrap: 'wrap',
    },
    chip: {
      margin: 2,
    },
    noLabel: {
      marginTop: theme.spacing(3),
    },
  }),
);

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250,
    },
  },
};

const types = [
  'Read A Book',
  'Meditate',
];


interface IProps extends RouteComponentProps {
  name?: string;
  ProjectAPI: ProjectAPI
}

export function ProjectPage(props: IProps) {
  const [state, setState] = useState({ metadata: {}} as Project);
  const id = (props.match.params as any)['id'];
  useEffect(() => {    
    const fetch = async () => {
      const prj = await props.ProjectAPI.get(id)
      setState(prj)
    }
    fetch();
  
  }, [props.ProjectAPI, id]);
  return (
    <div>{JSON.stringify(state)}</div>
  )
}