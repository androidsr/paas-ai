// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {menu} from '../models';
import {model} from '../models';
import {entity} from '../models';

export function AddMenus():Promise<menu.Menu>;

export function Close():Promise<void>;

export function GetConfig():Promise<model.HttpResult>;

export function GoToPage(arg1:string):Promise<void>;

export function SetConfig(arg1:entity.Config):Promise<model.HttpResult>;
