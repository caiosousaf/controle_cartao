package usuarios

import (
	"controle_cartao/infrastructure/cadastros/usuarios"
	"github.com/google/uuid"
)

// IUsuario define uma interface para os metodos de acesso a camada de dados
type IUsuario interface {
	CadastrarUsuario(req *usuarios.Usuario) error
	BuscarUsuarioLogin(nome *string) (*usuarios.Usuario, error)
	AtualizarSenhaUsuario(novaSenha, email *string, usuarioID *uuid.UUID) error
	BuscarUsuario(usuarioID *uuid.UUID) (res *usuarios.Usuario, err error)
}
