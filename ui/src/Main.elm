module Main exposing (main)

import Base64
import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (attribute, class, href, src, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as D
import Json.Encode as E
import Url exposing (Url)
import Url.Builder
import Url.Parser as Parser exposing ((</>), Parser)



-- MAIN


main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlRequest = UrlRequested
        , onUrlChange = UrlChanged
        }



-- MODEL


type FileTree
    = File String String
    | Directory String String (List FileTree) (List FileTree)


type alias Source =
    { name : String
    , path : String
    , content : String
    }


type alias Model =
    { id : Maybe String
    , source : Maybe Source
    , product : Maybe String
    , fileTree : Maybe FileTree
    , sidebarFlag : Bool
    , text : String
    , key : Nav.Key
    , url : Url
    }


type Msg
    = UrlChanged Url
    | UrlRequested Browser.UrlRequest
    | OpenSidebar
    | CloseSidebar
    | NewProjectRequest
    | NewProjectResult (Result Http.Error String)
    | FileTreeRequest String
    | FileTreeResult (Result Http.Error FileTree)
    | GetFileRequest String
    | GetFileResult (Result Http.Error Source)
    | CompileRequest
    | CompileResult (Result Http.Error String)
    | SaveRequest
    | Saved (Result Http.Error ())
    | SourceUpdate String
    | TextUpdate String
    | NewFileRequest
    | NewFileResult (Result Http.Error Source)



-- DECODER


fileTreeFileDecoder : D.Decoder FileTree
fileTreeFileDecoder =
    D.map2 File
        (D.field "name" D.string)
        (D.field "path" D.string)


fileTreeDecoder : D.Decoder FileTree
fileTreeDecoder =
    D.map4 Directory
        (D.field "name" D.string)
        (D.field "path" D.string)
        (D.field "childdirs" (D.list (D.lazy (\_ -> fileTreeDecoder))))
        (D.field "children" (D.list fileTreeFileDecoder))


sourceDecoder : D.Decoder Source
sourceDecoder =
    D.map3 Source
        (D.field "name" D.string)
        (D.field "path" D.string)
        (D.field "content" D.string)



-- INIT


projectParser : Parser (String -> a) a
projectParser =
    Parser.s "project" </> Parser.string


initWithID : String -> Url -> Nav.Key -> ( Model, Cmd Msg )
initWithID id url key =
    ( Model (Just id) Nothing Nothing Nothing False "" key url
    , fileTreeRequest id
    )


initWithNothing : Url -> Nav.Key -> ( Model, Cmd Msg )
initWithNothing url key =
    ( Model Nothing Nothing Nothing Nothing False "" key url
    , Cmd.none
    )


init : () -> Url -> Nav.Key -> ( Model, Cmd Msg )
init _ url key =
    case Parser.parse projectParser url of
        Just id ->
            initWithID id url key

        Nothing ->
            initWithNothing url key



-- UPDATE


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UrlRequested urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    {--( model, Nav.pushUrl model.key (Url.toString url) ) --}
                    ( model, Cmd.none )

                Browser.External href ->
                    ( model, Nav.load href )

        UrlChanged url ->
            case Parser.parse projectParser url of
                Just id ->
                    initWithID id model.url model.key

                Nothing ->
                    initWithNothing model.url model.key

        OpenSidebar ->
            ( { model | sidebarFlag = True }, Cmd.none )

        CloseSidebar ->
            ( { model | sidebarFlag = False }, Cmd.none )

        NewProjectRequest ->
            ( model, newProjectRequest )

        NewProjectResult (Result.Ok id) ->
            ( model, Nav.pushUrl model.key (Url.Builder.absolute [ "project", id ] []) )

        SourceUpdate newContent ->
            let
                newSource =
                    case model.source of
                        Just source ->
                            Just { source | content = newContent }

                        Nothing ->
                            Nothing
            in
            ( { model | source = newSource }, Cmd.none )

        TextUpdate newText ->
            ( { model | text = newText }, Cmd.none )

        SaveRequest ->
            case ( model.id, model.source ) of
                ( Just id, Just source ) ->
                    ( model, saveRequest id source )

                _ ->
                    ( model, Cmd.none )

        FileTreeRequest id ->
            ( model, fileTreeRequest id )

        FileTreeResult (Result.Ok fileTreeResult) ->
            ( { model | fileTree = Just fileTreeResult }, Cmd.none )

        FileTreeResult (Result.Err _) ->
            ( model, Cmd.none )

        GetFileRequest path ->
            case model.id of
                Just id ->
                    ( model, getFileRequest id path )

                Nothing ->
                    ( model, Cmd.none )

        GetFileResult (Result.Ok newSource) ->
            ( { model | source = Just newSource }, Cmd.none )

        GetFileResult (Result.Err _) ->
            ( model, Cmd.none )

        CompileRequest ->
            case ( model.id, model.source ) of
                ( Just id, Just source ) ->
                    ( { model | product = Nothing }, compileRequest id source )

                _ ->
                    ( model, Cmd.none )

        CompileResult (Result.Ok result) ->
            ( { model | product = Just result }, Cmd.none )

        CompileResult (Result.Err _) ->
            ( model, Cmd.none )

        NewFileRequest ->
            case model.id of
                Just id ->
                    ( { model | text = "" }, newFileRequest id model.text )

                Nothing ->
                    ( model, Cmd.none )

        NewFileResult (Result.Ok newFileSource) ->
            case model.id of
                Just id ->
                    ( { model | source = Just newFileSource }, fileTreeRequest id )

                Nothing ->
                    ( { model | source = Just newFileSource }, Cmd.none )

        NewFileResult (Result.Err _) ->
            ( model, Cmd.none )

        _ ->
            ( model, Cmd.none )



-- CMDS


newProjectRequest : Cmd Msg
newProjectRequest =
    Http.post
        { url = "/api/new-project"
        , body = Http.emptyBody
        , expect = Http.expectString NewProjectResult
        }


fileTreeRequest : String -> Cmd Msg
fileTreeRequest id =
    Http.get
        { url = Url.Builder.absolute [ "api", id, "list" ] []
        , expect = Http.expectJson FileTreeResult fileTreeDecoder
        }


getFileRequest : String -> String -> Cmd Msg
getFileRequest id path =
    Http.get
        { url = Url.Builder.absolute [ "api", id, "get" ] [ Url.Builder.string "path" path ]
        , expect = Http.expectJson GetFileResult sourceDecoder
        }


saveRequest : String -> Source -> Cmd Msg
saveRequest id source =
    Http.post
        { url = Url.Builder.absolute [ "api", id, "save" ] []
        , body =
            Http.jsonBody <|
                E.object
                    [ ( "path", E.string source.path )
                    , ( "data", E.string (Base64.encode source.content) )
                    ]
        , expect = Http.expectWhatever Saved
        }


newFileRequest : String -> String -> Cmd Msg
newFileRequest id path =
    Http.post
        { url = Url.Builder.absolute [ "api", id, "save" ] []
        , body =
            Http.jsonBody <|
                E.object
                    [ ( "path", E.string path )
                    , ( "data", E.string "" )
                    ]
        , expect =
            Http.expectWhatever
                (\r -> NewFileResult (Result.map (\_ -> Source "" path "") r))
        }


compileRequest : String -> Source -> Cmd Msg
compileRequest id source =
    Http.post
        { url = Url.Builder.absolute [ "api", id, "compile" ] []
        , body = Http.jsonBody <| E.object [ ( "path", E.string source.path ) ]
        , expect = Http.expectString CompileResult
        }



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Hello World"
    , body =
        [ div
            [ class "container"
            ]
            [ if model.sidebarFlag then
                div [ class "sidebar" ]
                    [ sidebarBackground
                        [ class "sidebarBackground"
                        , onClick CloseSidebar
                        ]
                        model
                    , fileTree
                        [ class "fileTree"
                        , onClick CloseSidebar
                        ]
                        model
                    ]

              else
                div [] []
            , div [ class "menu" ]
                [ button [ onClick OpenSidebar ] [ text "FILES" ]
                , button [ onClick SaveRequest ] [ text "SAVE" ]
                , button [ onClick CompileRequest ] [ text "COMPILE" ]
                , button [ onClick NewProjectRequest ] [ text "NEW PROJECT" ]
                , br [] []
                , div [ class "float-right" ]
                    [ span [] [ a [ href "//github.com/theoremoon/SATySFi-Online" ] [ text "Source" ] ]
                    ]
                , input [ type_ "text", value model.text, onInput TextUpdate ] []
                , button [ onClick NewFileRequest ] [ text "NEW FILE" ]
                , case model.source of
                    Just source ->
                        span [] [ text source.path ]

                    Nothing ->
                        span [] [ text "No file opened" ]
                ]
            , div [ class "main" ]
                [ satysfiEditor [ class "editor" ] model
                , productViewer [ class "viewer" ] model
                ]
            ]
        ]
    }


satysfiEditor : List (Attribute Msg) -> Model -> Html Msg
satysfiEditor attrs model =
    case model.source of
        Just source ->
            textarea
                (attrs
                    ++ [ value source.content
                       , onInput SourceUpdate
                       ]
                )
                []

        Nothing ->
            div attrs [ text "Please Load File" ]


productViewer : List (Attribute Msg) -> Model -> Html Msg
productViewer attrs model =
    case model.product of
        Just pdf ->
            embed
                (attrs
                    ++ [ src ("data:application/pdf;base64," ++ pdf)
                       ]
                )
                []

        Nothing ->
            div attrs [ text "Nothing there" ]


fileTreeImpl : FileTree -> Html Msg
fileTreeImpl filetree =
    case filetree of
        File name path ->
            li []
                [ span [ onClick (GetFileRequest path) ] [ text name ]
                ]

        Directory name _ dirs children ->
            li []
                [ text name
                , ul [] (List.map fileTreeImpl dirs)
                , ul [] (List.map fileTreeImpl children)
                ]


fileTree : List (Attribute Msg) -> Model -> Html Msg
fileTree attrs model =
    case model.fileTree of
        Just tree ->
            ul attrs [ fileTreeImpl tree ]

        Nothing ->
            ul attrs []


sidebarBackground : List (Attribute Msg) -> Model -> Html Msg
sidebarBackground attrs model =
    div attrs []
